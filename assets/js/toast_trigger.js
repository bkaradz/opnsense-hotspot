document.body.addEventListener('show-toast', function (evt) {
  const detail = evt.detail;

  // Create the toast container
  const toast = document.createElement('div');
  toast.setAttribute('data-tui-toast', '');
  toast.setAttribute('data-position', detail.position || 'bottom-right');
  toast.setAttribute('data-variant', detail.variant || 'default');
  toast.setAttribute('data-tui-toast-duration', detail.duration || '3000');

  // Base classes matching the templ component
  const baseClasses = "z-50 fixed pointer-events-auto p-4 w-full md:max-w-[420px] animate-in fade-in slide-in-from-bottom-4 duration-300";
  const positionClasses = "data-[position=bottom-right]:bottom-0 data-[position=bottom-right]:right-0 data-[position*=bottom]:slide-in-from-bottom-4";
  toast.className = `${baseClasses} ${positionClasses}`;

  // Inner content
  let iconHtml = '';
  if (detail.variant === 'success') {
    iconHtml = `<svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-circle-check text-green-500 mr-3 flex-shrink-0"><circle cx="12" cy="12" r="10"/><path d="m9 12 2 2 4-4"/></svg>`;
  } else if (detail.variant === 'error') {
    iconHtml = `<svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-circle-x text-red-500 mr-3 flex-shrink-0"><circle cx="12" cy="12" r="10"/><path d="m15 9-6 6"/><path d="m9 9 6 6"/></svg>`;
  }

  toast.innerHTML = `
        <div class="w-full bg-popover text-popover-foreground rounded-lg shadow-xs border pt-5 pb-4 px-4 flex items-center justify-center relative overflow-hidden group">
            <div class="absolute top-0 left-0 right-0 h-1 overflow-hidden">
                <div class="toast-progress h-full origin-left transition-transform ease-linear bg-gray-500" 
                     data-variant="${detail.variant || 'default'}"
                     style="background-color: ${detail.variant === 'success' ? '#22c55e' : (detail.variant === 'error' ? '#ef4444' : '')}">
                </div>
            </div>
            ${iconHtml}
            <span class="flex-1 min-w-0">
                <p class="text-sm font-semibold truncate">${detail.title || ''}</p>
                <p class="text-sm opacity-90 mt-1">${detail.description || ''}</p>
            </span>
            <button type="button" data-tui-toast-dismiss aria-label="Close" class="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-9 w-9">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-x opacity-75 hover:opacity-100"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
            </button>
        </div>
    `;

  document.body.appendChild(toast);
});
